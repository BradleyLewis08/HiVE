import React, { useEffect, useRef, useState } from "react";
import {
  Box,
  Button,
  Container,
  Flex,
  Grid,
  Heading,
  Text,
  Stack,
  Icon,
  useColorModeValue,
  VStack,
  HStack,
  Tag,
  BoxProps,
  Circle,
  useBreakpointValue,
} from "@chakra-ui/react";
import {
  motion,
  useScroll,
  useTransform,
  useSpring,
  useInView,
  Variants,
} from "framer-motion";
import {
  Code,
  Terminal,
  BookOpen,
  Server,
  Cpu,
  Users,
  LucideIcon,
  GitBranch,
  Workflow,
} from "lucide-react";

// Types
interface FeatureCardProps {
  icon: LucideIcon;
  title: string;
  description: string;
  index: number;
}

type MotionBoxProps = BoxProps & {
  children: React.ReactNode;
};

// Floating animation for background elements
const floatingAnimation: Variants = {
  initial: { y: 0 },
  animate: {
    y: [0, -20, 0],
    transition: {
      duration: 5,
      repeat: Infinity,
      repeatType: "reverse",
    },
  },
};

// Number animation component
const AnimatedNumber: React.FC<{ value: number; duration?: number }> = ({
  value,
  duration = 2,
}) => {
  const [count, setCount] = useState(0);
  const nodeRef = useRef<HTMLDivElement>(null);
  const inView = useInView(nodeRef);

  useEffect(() => {
    if (inView) {
      let start = 0;
      const end = value;
      const increment = end / (duration * 60);
      const timer = setInterval(() => {
        start += increment;
        if (start > end) {
          setCount(end);
          clearInterval(timer);
        } else {
          setCount(Math.floor(start));
        }
      }, 1000 / 60);
      return () => clearInterval(timer);
    }
  }, [value, duration, inView]);

  return <div ref={nodeRef}>{count.toLocaleString()}+</div>;
};

// Animated hexagon background
const HexagonBackground: React.FC = () => {
  return (
    <Box
      position="absolute"
      top={0}
      left={0}
      right={0}
      bottom={0}
      overflow="hidden"
      zIndex={0}
      opacity={0.1}
    >
      {Array.from({ length: 20 }).map((_, i) => (
        <motion.div
          key={i}
          style={{
            position: "absolute",
            left: `${Math.random() * 100}%`,
            top: `${Math.random() * 100}%`,
            width: "100px",
            height: "115px",
            background: `url("data:image/svg+xml,%3Csvg width='100' height='115' viewBox='0 0 100 115' fill='none' xmlns='http://www.w3.org/2000/svg'%3E%3Cpath d='M50 0L93.3013 25V75L50 100L6.69873 75V25L50 0Z' fill='%234299E1' fill-opacity='0.1'/%3E%3C/svg%3E")`,
          }}
          animate={{
            y: [0, -30, 0],
            rotate: [0, 360],
          }}
          transition={{
            duration: 20 + Math.random() * 10,
            repeat: Infinity,
            ease: "linear",
          }}
        />
      ))}
    </Box>
  );
};

// Enhanced feature card with hover effects and animations
const FeatureCard: React.FC<FeatureCardProps> = ({
  icon,
  title,
  description,
  index,
}) => {
  const bgColor = useColorModeValue("white", "gray.800");
  const borderColor = useColorModeValue("gray.100", "gray.700");
  const iconBgColor = useColorModeValue("blue.50", "blue.900");

  return (
    <motion.div
      initial={{ opacity: 0, y: 50 }}
      whileInView={{ opacity: 1, y: 0 }}
      viewport={{ once: true }}
      transition={{ duration: 0.5, delay: index * 0.1 }}
    >
      <Box
        p={8}
        bg={bgColor}
        rounded="xl"
        shadow="xl"
        borderWidth="1px"
        borderColor={borderColor}
        position="relative"
        overflow="hidden"
        _hover={{
          transform: "translateY(-8px)",
          shadow: "2xl",
          transition: "all 0.3s ease",
        }}
      >
        <motion.div
          whileHover={{ scale: 1.1 }}
          transition={{ type: "spring", stiffness: 300 }}
        >
          <Circle
            size="50px"
            bg={iconBgColor}
            mb={4}
            display="flex"
            alignItems="center"
            justifyContent="center"
          >
            <Icon as={icon} w={6} h={6} color="blue.500" />
          </Circle>
        </motion.div>

        <Text
          fontWeight="bold"
          fontSize="xl"
          mb={2}
          bgGradient="linear(to-r, blue.400, purple.500)"
          bgClip="text"
        >
          {title}
        </Text>
        <Text color={useColorModeValue("gray.600", "gray.400")}>
          {description}
        </Text>

        <Box
          position="absolute"
          right="-20px"
          bottom="-20px"
          width="100px"
          height="100px"
          opacity={0.1}
          bgGradient="radial(blue.400, transparent)"
          borderRadius="full"
        />
      </Box>
    </motion.div>
  );
};

// Stats section with animated counters
const StatsSection: React.FC = () => {
  const stats = [
    { label: "Active Users", value: 5000 },
    { label: "Universities", value: 150 },
    { label: "Countries", value: 45 },
    { label: "Development Hours Saved", value: 25000 },
  ];

  return (
    <Grid
      templateColumns={{ base: "repeat(2, 1fr)", md: "repeat(4, 1fr)" }}
      gap={8}
      my={20}
    >
      {stats.map(({ label, value }, index) => (
        <motion.div
          key={label}
          initial={{ opacity: 0, scale: 0.5 }}
          whileInView={{ opacity: 1, scale: 1 }}
          viewport={{ once: true }}
          transition={{ duration: 0.5, delay: index * 0.1 }}
        >
          <VStack p={6} bg="white" rounded="xl" shadow="lg" spacing={2}>
            <Text fontSize="4xl" fontWeight="bold" color="blue.500">
              <AnimatedNumber value={value} />
            </Text>
            <Text color="gray.600">{label}</Text>
          </VStack>
        </motion.div>
      ))}
    </Grid>
  );
};

// Main component
const Landing: React.FC = () => {
  const bgColor = useColorModeValue("gray.50", "gray.900");
  const { scrollYProgress } = useScroll();
  const isMobile = useBreakpointValue({ base: true, md: false });

  const scaleProgess = useTransform(scrollYProgress, [0, 1], [1, 0.8]);
  const opacityProgess = useTransform(scrollYProgress, [0, 1], [1, 0.3]);

  const springConfig = { stiffness: 100, damping: 30, restDelta: 0.001 };
  const scaleX = useSpring(scrollYProgress, springConfig);

  return (
    <Box bg={bgColor} minH="100vh" overflowX="hidden">
      <HexagonBackground />

      {/* Progress bar */}
      <motion.div
        style={{
          position: "fixed",
          top: 0,
          left: 0,
          right: 0,
          height: "4px",
          background: "linear-gradient(to right, #4299E1, #805AD5)",
          transformOrigin: "0%",
          scaleX,
          zIndex: 100,
        }}
      />

      <Container maxW="7xl" pt={20} pb={20} position="relative">
        {/* Hero Section with Parallax */}
        <motion.div
          style={{
            scale: scaleProgess,
            opacity: opacityProgess,
          }}
        >
          <Stack spacing={12} align="center" position="relative">
            <motion.div
              initial={{ opacity: 0, y: -50 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.8, ease: "easeOut" }}
            >
              <Heading
                fontSize={{ base: "4xl", md: "6xl" }}
                fontWeight="bold"
                lineHeight="shorter"
                mb={6}
                textAlign="center"
              >
                Welcome to{" "}
                <motion.span
                  style={{
                    display: "inline-block",
                    background: "linear-gradient(to right, #4299E1, #805AD5)",
                    WebkitBackgroundClip: "text",
                    WebkitTextFillColor: "transparent",
                  }}
                  animate={{ scale: [1, 1.2, 1] }}
                  transition={{
                    duration: 2,
                    repeat: Infinity,
                    repeatType: "reverse",
                  }}
                >
                  HiVE
                </motion.span>
              </Heading>
            </motion.div>

            {/* Floating elements */}
            {!isMobile && (
              <>
                <motion.div
                  variants={floatingAnimation}
                  initial="initial"
                  animate="animate"
                  style={{
                    position: "absolute",
                    left: "-10%",
                    top: "20%",
                  }}
                >
                  <Icon as={GitBranch} w={12} h={12} color="blue.200" />
                </motion.div>
                <motion.div
                  variants={floatingAnimation}
                  initial="initial"
                  animate="animate"
                  style={{
                    position: "absolute",
                    right: "-5%",
                    top: "40%",
                  }}
                >
                  <Icon as={Workflow} w={12} h={12} color="purple.200" />
                </motion.div>
              </>
            )}

            {/* Tags with staggered animation */}
            <motion.div
              initial="hidden"
              animate="visible"
              variants={{
                hidden: { opacity: 0 },
                visible: {
                  opacity: 1,
                  transition: {
                    staggerChildren: 0.1,
                  },
                },
              }}
            >
              <HStack spacing={4} justify="center" wrap="wrap">
                {[
                  "Containerized",
                  "Cloud-Native",
                  "Scalable",
                  "Secure",
                  "Collaborative",
                ].map((tag, index) => (
                  <motion.div
                    key={tag}
                    variants={{
                      hidden: { opacity: 0, x: -20 },
                      visible: { opacity: 1, x: 0 },
                    }}
                  >
                    <Tag
                      size="lg"
                      colorScheme={
                        ["blue", "purple", "green", "red", "orange"][index]
                      }
                      rounded="full"
                      px={6}
                      py={2}
                      fontSize="md"
                    >
                      {tag}
                    </Tag>
                  </motion.div>
                ))}
              </HStack>
            </motion.div>
          </Stack>
        </motion.div>

        {/* Stats Section */}
        <StatsSection />

        {/* Features Grid with Staggered Animation */}
        <Grid
          templateColumns={{
            base: "1fr",
            md: "repeat(2, 1fr)",
            lg: "repeat(3, 1fr)",
          }}
          gap={8}
          mt={20}
        >
          {[
            {
              icon: Terminal,
              title: "Standardized Environments",
              description:
                "Ensure consistency across different computing environments with containerized development setups.",
            },
            {
              icon: Server,
              title: "Cloud-Agnostic",
              description:
                "Deploy on various cloud platforms or on-premises infrastructure with our flexible architecture.",
            },
            {
              icon: Code,
              title: "VSCode Integration",
              description:
                "Seamlessly connect to your development environment through our VSCode extension.",
            },
            {
              icon: BookOpen,
              title: "Educational Focus",
              description:
                "Purpose-built for computer science education with features designed for both students and educators.",
            },
            {
              icon: Cpu,
              title: "Resource Optimization",
              description:
                "Efficient resource allocation and scaling to accommodate varying class sizes and computational demands.",
            },
            {
              icon: Users,
              title: "Collaborative Learning",
              description:
                "Enable seamless collaboration between students and instructors with shared environments.",
            },
          ].map((feature, index) => (
            <FeatureCard key={index} {...feature} index={index} />
          ))}
        </Grid>

        {/* CTA Section with Scale Animation */}
        <motion.div
          initial={{ opacity: 0, scale: 0.8 }}
          whileInView={{ opacity: 1, scale: 1 }}
          viewport={{ once: true }}
          transition={{ duration: 0.5 }}
        >
          <VStack spacing={6} mt={20} textAlign="center">
            <Heading size="lg">Ready to transform CS education?</Heading>
            <motion.div whileHover={{ scale: 1.05 }} whileTap={{ scale: 0.95 }}>
              <Button
                size="lg"
                colorScheme="blue"
                rounded="full"
                px={8}
                shadow="lg"
                bgGradient="linear(to-r, blue.400, purple.500)"
                _hover={{
                  bgGradient: "linear(to-r, blue.500, purple.600)",
                  transform: "translateY(-2px)",
                  shadow: "xl",
                }}
              >
                Request Demo
              </Button>
            </motion.div>
          </VStack>
        </motion.div>
      </Container>
    </Box>
  );
};

export default Landing;

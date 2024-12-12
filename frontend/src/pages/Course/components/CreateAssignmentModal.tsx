import {
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalFooter,
  ModalBody,
  ModalCloseButton,
  FormControl,
  FormLabel,
  Input,
  Textarea,
  VStack,
  Button,
} from "@chakra-ui/react";
import { AssignmentFormData } from "../Course";

interface CreateAssignmentModalProps {
  isOpen: boolean;
  onClose: () => void;
  onInputChange: (
    e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>
  ) => void;
  formData: AssignmentFormData;
  handleCreateAssignment: () => void;
}

export default function CreateAssignmentModal({
  isOpen,
  onClose,
  onInputChange,
  formData,
  handleCreateAssignment,
}: CreateAssignmentModalProps) {
  return (
    <Modal isOpen={isOpen} onClose={onClose}>
      <ModalOverlay />
      <ModalContent>
        <ModalHeader>Create New Assignment</ModalHeader>
        <ModalCloseButton />
        <ModalBody>
          <VStack spacing={4}>
            <FormControl isRequired>
              <FormLabel>Title</FormLabel>
              <Input
                name="title"
                placeholder="Enter assignment title"
                value={formData.title}
                onChange={onInputChange}
              />
            </FormControl>
            <FormControl isRequired>
              <FormLabel>Description</FormLabel>
              <Textarea
                name="description"
                placeholder="Enter assignment description"
                value={formData.description}
                onChange={onInputChange}
              />
            </FormControl>
            <FormControl isRequired>
              <FormLabel>Image URL</FormLabel>
              <Input
                name="imageURL"
                placeholder="Enter image URL"
                value={formData.imageURL}
                onChange={onInputChange}
              />
            </FormControl>
          </VStack>
        </ModalBody>

        <ModalFooter>
          <Button variant="ghost" mr={3} onClick={onClose}>
            Cancel
          </Button>
          <Button
            colorScheme="blue"
            onClick={handleCreateAssignment}
            isDisabled={
              !formData.title || !formData.description || !formData.imageURL
            }
          >
            Create
          </Button>
        </ModalFooter>
      </ModalContent>
    </Modal>
  );
}

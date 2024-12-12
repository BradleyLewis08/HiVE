import { useEffect } from "react";
import { useSearchParams, useNavigate } from "react-router-dom";
import { useAuth } from "../../providers/AuthProvider";
import CenteredLoader from "../Loading/CenteredLoader";

export default function Redirect() {
  let navigate = useNavigate();
  const { initializeAuth } = useAuth();
  let [searchParams] = useSearchParams();

  // Handle token interception
  useEffect(() => {
    async function handleRedirect() {
      const token = searchParams.get("token");
      if (token == null) {
        navigate("/login");
        return;
      }

      try {
        await initializeAuth(token);
        navigate("/home");
      } catch (error) {
        console.error("Error initializing auth:", error);
        navigate("/login");
      }
    }

    handleRedirect();
  }, [initializeAuth, navigate, searchParams]);

  return <CenteredLoader />;
}

import {
  createBrowserRouter,
  RouterProvider,
  Navigate,
  useLocation,
} from "react-router-dom";
import Redirect from "../pages/Redirect";
import { useAuth } from "../providers/AuthProvider";
import Landing from "../pages/Landing/Landing";
import Home from "../pages/Home";
import Course from "../pages/Course/Course";

const ProtectedRoute = ({ children }: { children: React.ReactNode }) => {
  const { isAuthenticated } = useAuth();

  const location = useLocation();

  if (!isAuthenticated) {
    return <Navigate to="/landing" state={{ from: location }} />;
  }

  return <>{children}</>;
};

const PublicRoute = ({ children }: { children: React.ReactNode }) => {
  const { isAuthenticated } = useAuth();

  if (isAuthenticated) {
    return <Navigate to="/home" />;
  }

  return <>{children}</>;
};

const routes = [
  {
    path: "/landing",
    element: (
      <PublicRoute>
        <Landing />
      </PublicRoute>
    ),
  },
  {
    path: "/redirect",
    element: <Redirect />,
  },
  {
    path: "/home",
    element: (
      <ProtectedRoute>
        <Home />
      </ProtectedRoute>
    ),
  },
  {
    path: "/course/:courseCode",
    element: (
      <ProtectedRoute>
        <Course />
      </ProtectedRoute>
    ),
  },
  {
    path: "*",
    element: <Navigate to="/landing" />,
  },
];

const browserRouter = createBrowserRouter(routes);

export default function Router() {
  return <RouterProvider router={browserRouter} />;
}

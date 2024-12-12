import Router from "./routes/router";
import Providers from "./providers";
import StickyNavbar from "./components/nav/Navbar";

export const App = () => (
  <Providers>
    <StickyNavbar />
    <Router />
  </Providers>
);

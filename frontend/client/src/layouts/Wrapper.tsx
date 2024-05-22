import "./styles.css";

import { FirebaseProvider } from "@/hooks/firebase";
import { FrontendServiceProvider } from "@/hooks/rpc";

export default function Wrapper({ children }: { children: React.ReactNode }) {
  return (
    <FirebaseProvider>
      <FrontendServiceProvider>{children}</FrontendServiceProvider>
    </FirebaseProvider>
  );
}

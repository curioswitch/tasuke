import "./styles.css";

import { FirebaseProvider } from "@/hooks/firebase";
import { FrontendServiceProvider } from "@/hooks/rpc";
import { NextUIProvider } from "@nextui-org/system";

export default function Wrapper({ children }: { children: React.ReactNode }) {
  return (
    <FirebaseProvider>
      <FrontendServiceProvider>
        <NextUIProvider>{children}</NextUIProvider>
      </FrontendServiceProvider>
    </FirebaseProvider>
  );
}

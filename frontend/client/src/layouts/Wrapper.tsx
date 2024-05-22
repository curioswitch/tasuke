import "./styles.css";

import { NextUIProvider } from "@nextui-org/system";
import { navigate } from "vike/client/router";

import { FirebaseProvider } from "@/hooks/firebase";
import { FrontendServiceProvider } from "@/hooks/rpc";

export default function Wrapper({ children }: { children: React.ReactNode }) {
  return (
    <FirebaseProvider>
      <FrontendServiceProvider>
        <NextUIProvider navigate={navigate}>{children}</NextUIProvider>
      </FrontendServiceProvider>
    </FirebaseProvider>
  );
}

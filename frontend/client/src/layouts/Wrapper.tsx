import "./styles.css";

import { NextUIProvider } from "@nextui-org/system";
import { usePageContext } from "vike-react/usePageContext";
import { navigate } from "vike/client/router";

import { FirebaseProvider, useFirebase } from "@/hooks/firebase";
import { FrontendServiceProvider } from "@/hooks/rpc";

const authorizedPaths = ["/settings"];

function Authorizer({ children }: { children: React.ReactNode }) {
  const firebase = useFirebase();
  const pageCtx = usePageContext();

  if (!authorizedPaths.includes(pageCtx.urlPathname)) {
    return <>{children}</>;
  }

  if (!firebase?.userResolved) {
    return;
  }

  if (!firebase.user) {
    navigate("/");
    return;
  }

  return <>{children}</>;
}

export default function Wrapper({ children }: { children: React.ReactNode }) {
  return (
    <FirebaseProvider>
      <FrontendServiceProvider>
        <NextUIProvider navigate={navigate}>
          <Authorizer>{children}</Authorizer>
        </NextUIProvider>
      </FrontendServiceProvider>
    </FirebaseProvider>
  );
}

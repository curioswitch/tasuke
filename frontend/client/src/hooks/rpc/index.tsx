import type { Interceptor } from "@connectrpc/connect";
import { TransportProvider, useQuery } from "@connectrpc/connect-query";
import { createConnectTransport } from "@connectrpc/connect-web";
import { type User, getUser } from "@tasuke/frontendapi";
import type { User as FirebaseUser } from "firebase/auth";
import { useMemo } from "react";
import { useFirebaseUser } from "../firebase";

function createFirebaseAuthInterceptor(user: FirebaseUser): Interceptor {
  return (next) => async (request) => {
    const idToken = await user.getIdToken();
    request.header.set("authorization", `Bearer ${idToken}`);
    return next(request);
  };
}

export function useUser(): User | undefined {
  const fbUser = useFirebaseUser();

  const { data } = useQuery(getUser, undefined, {
    enabled: !!fbUser,
  });

  return data?.user;
}

export function FrontendServiceProvider({
  children,
}: { children: React.ReactNode }) {
  const fbUser = useFirebaseUser();

  const transport = useMemo(() => {
    const interceptors = fbUser ? [createFirebaseAuthInterceptor(fbUser)] : [];
    return createConnectTransport({ baseUrl: "/", interceptors });
  }, [fbUser]);

  return (
    <TransportProvider transport={transport}>{children}</TransportProvider>
  );
}

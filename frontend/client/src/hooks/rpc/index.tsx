import type { Interceptor } from "@connectrpc/connect";
import { TransportProvider, useQuery } from "@connectrpc/connect-query";
import { createConnectTransport } from "@connectrpc/connect-web";
import { QueryClient, QueryClientProvider } from "@tanstack/react-query";
import { type User, getUser } from "@tasuke/frontendapi";
import type { User as FirebaseUser } from "firebase/auth";
import { useMemo } from "react";

import { useFirebase, useFirebaseUser } from "@/hooks/firebase";

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
  const fbUser = useFirebase()?.user;

  const transport = useMemo(() => {
    const interceptors = fbUser ? [createFirebaseAuthInterceptor(fbUser)] : [];
    return createConnectTransport({ baseUrl: "/", interceptors });
  }, [fbUser]);

  // Most documentation puts the queryClient at the package level as it's the
  // easiest way to guarantee the cache is preserved regardless of rendering
  // loop mistakes, but we get to 1) avoid RPCs until firebase auth is resolved
  // and 2) easily invalidate the query cache on logout, potentially logging-in
  // as a different user, so we use useMemo here. The documentation says useMemo
  // should only be used as an optimization and it's not guaranteed to always
  // be stable even if the inputs don't change. While in some sense an optimization,
  // it is also important to not send unnecessary backend calls so we do want this
  // to be stable. In the unlikely event that useMemo behavior changes to be less
  // strict in the future, we can write our own stable hook instead.
  const queryClient = useMemo(() => {
    return new QueryClient({
      defaultOptions: {
        queries: {
          enabled: !!fbUser,
        },
      },
    });
  }, [fbUser]);

  return (
    <QueryClientProvider client={queryClient}>
      <TransportProvider transport={transport}>{children}</TransportProvider>
    </QueryClientProvider>
  );
}

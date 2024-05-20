import "./styles.css";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

import { FirebaseProvider } from "@/hooks/firebase";
import { FrontendServiceProvider } from "@/hooks/rpc";

const queryClient = new QueryClient();

export default function Wrapper({ children }: { children: React.ReactNode }) {
  return (
    <FirebaseProvider>
      <QueryClientProvider client={queryClient}>
        <FrontendServiceProvider>{children}</FrontendServiceProvider>
      </QueryClientProvider>
    </FirebaseProvider>
  );
}

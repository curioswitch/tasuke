import "./styles.css";

import { QueryClient, QueryClientProvider } from "@tanstack/react-query";

import { FirebaseProvider } from "@/hooks/firebase";
import { FrontendServiceProvider } from "@/hooks/rpc";

const queryClient = new QueryClient();

export default function Wrapper({ children }: { children: React.ReactNode }) {
  return (
    <QueryClientProvider client={queryClient}>
      <FirebaseProvider>
        <FrontendServiceProvider>{children}</FrontendServiceProvider>
      </FirebaseProvider>
    </QueryClientProvider>
  );
}

import { ChakraProvider } from "@chakra-ui/react";

export default function Wrapper({ children }: { children: React.ReactNode }) {
  return <ChakraProvider>{children}</ChakraProvider>;
}

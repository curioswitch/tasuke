import { FirebaseProvider } from "../hooks/firebase";

export default function Wrapper({ children }: { children: React.ReactNode }) {
  return <FirebaseProvider>{children}</FirebaseProvider>;
}

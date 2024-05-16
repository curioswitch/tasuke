import { Button } from "@/components/ui/button";
import { useAuth } from "@/hooks/firebase";
import { GithubAuthProvider, signInWithPopup } from "firebase/auth";
import { useCallback } from "react";

export default function Page() {
  const auth = useAuth();

  const onSignUpClick = useCallback(() => {
    if (!auth) {
      return;
    }

    signInWithPopup(auth, new GithubAuthProvider());
  }, [auth]);

  return (
    <div className="col-span-4 md:col-span-8 lg:col-span-12">
      <h1>Tasuke</h1>
      <Button onClick={onSignUpClick}>Sign up with GitHub</Button>
    </div>
  );
}

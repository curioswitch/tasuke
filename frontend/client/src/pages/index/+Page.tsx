import { Button } from "@nextui-org/button";
import { getApp } from "firebase/app";
import { GithubAuthProvider, getAuth, signInWithPopup } from "firebase/auth";
import { useCallback } from "react";
import { navigate } from "vike/client/router";

import { useFirebase } from "@/hooks/firebase";

export default function Page() {
  const fbState = useFirebase();

  const onSignUpClick = useCallback(async () => {
    await signInWithPopup(getAuth(getApp()), new GithubAuthProvider());
    navigate("/settings");
  }, []);

  return (
    <>
      <div className="col-span-4 md:col-span-8 lg:col-span-12">
        <h1>Tasuke</h1>
        <p className="lead">
          tasuke is short for the Japanese 手助け (tedasuke), meaning to give a
          helping hand. It aims to support OSS developers in giving each other a
          helping hand by connecting code reviewers to PRs in otherwise
          unrelated OSS repositories.
        </p>
        <p>
          Everyone knows the{" "}
          <a
            href="https://xkcd.com/2347/"
            className="underline"
            target="_blank"
            rel="noreferrer noopener"
          >
            lone maintainer
          </a>{" "}
          issue in OSS - let's see if we can make it at least a little less
          lonely.
        </p>
        <p>
          If you are interested in helping with code reviews, create an account!
        </p>
        {fbState?.userResolved && !fbState.user ? (
          <Button
            className="bg-primary-500 text-content1"
            onClick={onSignUpClick}
          >
            Sign up with GitHub
          </Button>
        ) : null}
      </div>
    </>
  );
}

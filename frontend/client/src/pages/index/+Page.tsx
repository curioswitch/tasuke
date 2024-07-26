import { Button } from "@nextui-org/button";
import { Image } from "@nextui-org/image";
import { Link } from "@nextui-org/link";
import { getApp } from "firebase/app";
import { GithubAuthProvider, getAuth, signInWithPopup } from "firebase/auth";
import { useCallback, useRef } from "react";
import { BiLogoGithub as LogoGithub } from "react-icons/bi";
import { navigate } from "vike/client/router";

import { useFirebase } from "@/hooks/firebase";

import assistanceImg from "./static/assistance.svg";
import bannerImg from "./static/banner.svg";
import handImg from "./static/hand.svg";
import handshakeImg from "./static/handshake.svg";
import laptopPhoneImg from "./static/laptop-phone.svg";
import maintainerImg from "./static/maintainer.svg";
import reviewerImg from "./static/reviewer.svg";

function getBotInstallLink() {
  if (import.meta.env.PUBLIC_ENV__FIREBASE_APP === "tasuke-dev") {
    return "https://github.com/apps/tasuke-alpha";
  }
  if (import.meta.env.PUBLIC_ENV__FIREBASE_APP === "tasuke-prod") {
    return "https://github.com/apps/tasuke-bot";
  }
  throw new Error("PUBLIC_ENV__FIREBASE_APP must be configured");
}

export default function Page() {
  const fbState = useFirebase();

  const onSignUpClick = useCallback(async () => {
    await signInWithPopup(getAuth(getApp()), new GithubAuthProvider());
    navigate("/settings");
  }, []);

  const forReviewerRef = useRef<HTMLDivElement>(null);
  const scrollToForReviewer = useCallback(() => {
    forReviewerRef.current?.scrollIntoView({ behavior: "smooth" });
  }, []);

  const forMaintainerRef = useRef<HTMLDivElement>(null);
  const scrollToForMaintainer = useCallback(() => {
    forMaintainerRef.current?.scrollIntoView({ behavior: "smooth" });
  }, []);
  const contributingRef = useRef<HTMLDivElement>(null);
  const scrollToContributing = useCallback(() => {
    contributingRef.current?.scrollIntoView({ behavior: "smooth" });
  }, []);

  return (
    <>
      <div className="col-span-4 md:col-span-8 lg:col-span-12">
        <div
          style={{ backgroundImage: `url(${bannerImg})` }}
          className="flex flex-col-reverse md:flex-row md:gap-10 p-20"
        >
          <div className="basis-1/2" />
          <div className="basis-1/2">
            <h2 className="text-white drop-shadow-[0_1.2px_1.2px_rgba(0,0,0,0.8)]">
              A mutual assistance service for OSS developers.
            </h2>
          </div>
        </div>
        <div className="flex flex-col-reverse md:flex-row md:gap-10 p-4 md:p-20">
          <div className="basis-1/2">
            <h2>ABOUT TASUKE</h2>
            <p>
              tasuke is short for the Japanese 手助け (tedasuke), meaning to
              give a helping hand. It aims to support OSS developers in giving
              each other a helping hand by connecting code reviewers to PRs in
              otherwise unrelated OSS repositories.
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
          </div>
          <div className="basis-1/2">
            <Image classNames={{ wrapper: "mx-auto" }} src={assistanceImg} />
          </div>
        </div>

        <div className="bg-primary-50 p-10 pt-20">
          <div className="text-center text-primary-400">What we serve</div>
          <h2 className="text-center m-0 mb-10">Getting Started</h2>

          <div className="flex flex-col md:flex-row gap-10 text-center lg:px-20">
            <button
              type="button"
              className="basis-1/3 bg-white rounded-medium border-1 py-10"
              onClick={scrollToForReviewer}
            >
              <Image classNames={{ wrapper: "mx-auto" }} src={handshakeImg} />
              <h4>For code reviewers</h4>
            </button>
            <button
              type="button"
              className="basis-1/3 bg-white rounded-medium border-1 py-10"
              onClick={scrollToForMaintainer}
            >
              <Image classNames={{ wrapper: "mx-auto" }} src={handImg} />
              <h4>For maintainers</h4>
            </button>
            <button
              type="button"
              className="basis-1/3 bg-white rounded-medium border-1 py-10"
              onClick={scrollToContributing}
            >
              <Image classNames={{ wrapper: "mx-auto" }} src={laptopPhoneImg} />
              <h4>Contributing to tasuke</h4>
            </button>
          </div>
        </div>

        <div
          ref={forReviewerRef}
          className="flex flex-col md:flex-row md:gap-20 mt-20 items-center p-4"
        >
          <div className="basis-1/2">
            <Image classNames={{ wrapper: "mx-auto" }} src={reviewerImg} />
          </div>
          <div className="basis-1/2">
            <h2 className="mt-0">For code reviewer</h2>
            {fbState?.userResolved && !fbState.user ? (
              <>
                <p>
                  If you are interested in helping with code reviews, create an
                  account and submit your preferences for reviewing.
                </p>
                <Button
                  className="bg-primary-400 text-content1"
                  onClick={onSignUpClick}
                  startContent={<LogoGithub className="size-6" />}
                >
                  Sign up with GitHub
                </Button>
              </>
            ) : (
              <>
                <p>
                  If you are interested in helping with code reviews, submit or
                  edit your preferences for reviewing.
                </p>
                <p>
                  <a href="/settings">Edit your preferences</a>
                </p>
              </>
            )}
          </div>
        </div>

        <div
          ref={forMaintainerRef}
          className="flex flex-col md:flex-row md:gap-20 mt-20 items-center p-4"
        >
          <div className="basis-1/2">
            <Image classNames={{ wrapper: "mx-auto" }} src={maintainerImg} />
          </div>
          <div className="basis-1/2">
            <h2 className="mt-0">For maintainers</h2>
            <p className="lead">
              If you have a repository you could use help with on code review,
              register the Tasuke GitHub app.
            </p>
            <p>
              Currently you must also register as a code reviewer to use the
              app. As review capacity increases and we understand abuse
              potential better, we may remove this requirement.
            </p>
            <p>
              <Button
                as={Link}
                className="bg-primary-400 text-content1"
                href={getBotInstallLink()}
                target="_blank"
                rel="noreferrer noopener"
                startContent={<LogoGithub className="size-6" />}
              >
                Install the app
              </Button>
            </p>
            <p>
              With the app installed on a repository, create a pull request in
              it and request a review by adding a comment starting with{" "}
              <code>/tasuke</code>, for example:
            </p>
            <p>
              <code>/tasuke Can you help me with a review, please?</code>
            </p>
            <p>
              Currently, the only supported command is to ask for a review, so
              the text content is not used. In the future, we may use it for
              executing different commands.
            </p>
          </div>
        </div>

        <div ref={contributingRef} className="bg-primary-50 p-4 md:p-20">
          <div className="flex items-center flex-col md:flex-row md:gap-20 bg-primary-400 text-content1 p-10 rounded-small">
            <div className="basis-2/3">
              <p>tasuke is OSS for OSS</p>
              <h1 className="text-content1">Contributing to tasuke</h1>
              <p>
                If you are interested, feel free to file issues or pull requests
                as with any OSS project. Let's work together to find what can
                make the service most useful for OSS developers.
              </p>
            </div>
            <div className="basis-1/3">
              <Button
                as={Link}
                className="bg-white text-primary-400"
                href="https://github.com/curioswitch/tasuke"
                target="_blank"
                rel="noreferrer noopener"
                showAnchorIcon={true}
              >
                Check out the repo
              </Button>
            </div>
          </div>
        </div>
      </div>
    </>
  );
}

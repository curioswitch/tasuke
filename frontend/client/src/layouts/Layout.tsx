import type React from "react";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { useFirebase, useFirebaseUser } from "@/hooks/firebase";
import { type User as FirebaseUser, signOut } from "firebase/auth";
import { LogOut, User } from "lucide-react";
import { forwardRef, useCallback } from "react";

const UserAvatar = forwardRef(
  ({ fbUser, ...props }: { fbUser?: FirebaseUser }, ref) => {
    return (
      <Avatar ref={ref} {...props}>
        <AvatarImage src={fbUser?.photoURL} />
        <AvatarFallback />
      </Avatar>
    );
  },
);

export default function Layout({
  children,
}: {
  children: React.ReactNode;
}) {
  const firebase = useFirebase();
  const fbUser = useFirebaseUser();

  const onLogOutClick = useCallback(() => {
    if (!firebase || !fbUser) {
      return;
    }

    signOut(firebase.auth);
  }, [firebase, fbUser]);

  return (
    <>
      <section>
        <div className="container px-6 m-auto">
          <div className="grid grid-cols-4 gap-6 md:grid-cols-8 lg:grid-cols-12">
            <div className="flex justify-end col-span-4 md:col-span-8 lg:col-span-12">
              {fbUser ? (
                <DropdownMenu>
                  <DropdownMenuTrigger asChild className="cursor-pointer">
                    <UserAvatar fbUser={fbUser} />
                  </DropdownMenuTrigger>
                  <DropdownMenuContent>
                    <DropdownMenuItem>
                      <User className="mr-2 h-4 w-4" />
                      <span>Profile</span>
                    </DropdownMenuItem>
                    <DropdownMenuItem onClick={onLogOutClick}>
                      <LogOut className="mr-2 h-4 w-4" />
                      <span>Log out</span>
                    </DropdownMenuItem>
                  </DropdownMenuContent>
                </DropdownMenu>
              ) : (
                <UserAvatar fbUser={fbUser} />
              )}
            </div>
            {children}
          </div>
        </div>
      </section>
    </>
  );
}

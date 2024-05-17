import type React from "react";

import { Avatar, AvatarFallback, AvatarImage } from "@/components/ui/avatar";
import { useFirebaseUser } from "@/hooks/firebase";

export default function Layout({
  children,
}: {
  children: React.ReactNode;
}) {
  const fbUser = useFirebaseUser();

  return (
    <>
      <section>
        <div className="container px-6 m-auto">
          <div className="grid grid-cols-4 gap-6 md:grid-cols-8 lg:grid-cols-12">
            <div className="flex justify-end col-span-4 md:col-span-8 lg:col-span-12">
              <Avatar>
                <AvatarImage src={fbUser?.photoURL} />
                <AvatarFallback />
              </Avatar>
            </div>
            {children}
          </div>
        </div>
      </section>
    </>
  );
}

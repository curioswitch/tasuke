import { useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@nextui-org/button";
import { Input } from "@nextui-org/input";
import { useQueryClient } from "@tanstack/react-query";
import { getUser } from "@tasuke/frontendapi";
import { useCallback } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { H1 } from "@/components/ui/typography";

const formSchema = z
  .object({
    programmingLanguageIds: z.array(z.number().int()),
    maxOpenReviews: z.number().int().nonnegative(),
  })
  .required();

export default function Page() {
  // Workaround https://github.com/connectrpc/connect-query-es/pull/369
  const queryClient = useQueryClient();
  const result = useQuery(getUser, undefined, {
    enabled: queryClient.getDefaultOptions().queries?.enabled,
  });

  const { data, isSuccess } = result;

  const user = data?.user;

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      programmingLanguageIds: user?.programmingLanguageIds ?? [],
      maxOpenReviews: user?.maxOpenReviews ?? 1,
    },
  });

  const onFormSubmit = useCallback((values: z.infer<typeof formSchema>) => {
    console.log(values);
  }, []);

  if (!isSuccess) {
    // TODO: Better loading screen.
    return <div>Loading...</div>;
  }

  return (
    <>
      <div className="col-span-4 md:col-span-8 lg:col-span-12">
        <H1>Profile</H1>
      </div>
      <div className="col-span-4 md:col-span-8 lg:col-span-12">
        <form onSubmit={handleSubmit(onFormSubmit)} className="space-y-8">
          <Input
            {...register("maxOpenReviews", { valueAsNumber: true })}
            type="number"
            label="Maximum open reviews"
            description="Should be at least 1 to help with reviews, but feel free to set to 0 on breaks."
            isInvalid={!!errors.maxOpenReviews}
            errorMessage={errors.maxOpenReviews?.message}
          />
          <Button type="submit" className="bg-primary-500 text-content1">
            Submit
          </Button>
        </form>
      </div>
    </>
  );
}

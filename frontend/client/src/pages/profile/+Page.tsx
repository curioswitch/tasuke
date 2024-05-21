import { useQuery } from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { useQueryClient } from "@tanstack/react-query";
import { getUser } from "@tasuke/frontendapi";
import { useCallback } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

import { Button } from "@/components/ui/button";
import {
  Form,
  FormControl,
  FormDescription,
  FormField,
  FormItem,
  FormLabel,
  FormMessage,
} from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { H1 } from "@/components/ui/typography";

const formSchema = z
  .object({
    programmingLanguageIds: z.array(z.number().int()),
    maxOpenReviews: z.preprocess((val) => {
      return val !== "" ? Number(val) : val;
    }, z.number().int()),
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

  const form = useForm<z.infer<typeof formSchema>>({
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
        <Form {...form}>
          <form
            onSubmit={form.handleSubmit(onFormSubmit)}
            className="space-y-8"
          >
            <FormField
              control={form.control}
              name="maxOpenReviews"
              render={({ field }) => (
                <FormItem>
                  <FormLabel>Maximum open reviews</FormLabel>
                  <FormControl>
                    <Input type="number" {...field} />
                  </FormControl>
                  <FormDescription>
                    Should be at least 1 to help with reviews, but feel free to
                    set to 0 on breaks.
                  </FormDescription>
                  <FormMessage />
                </FormItem>
              )}
            />
            <Button type="submit">Submit</Button>
          </form>
        </Form>
      </div>
    </>
  );
}

import {
  createConnectQueryKey,
  useMutation,
  useQuery,
} from "@connectrpc/connect-query";
import { zodResolver } from "@hookform/resolvers/zod";
import { Button } from "@nextui-org/button";
import { Input } from "@nextui-org/input";
import { useQueryClient } from "@tanstack/react-query";
import {
  SaveUserRequest,
  type User,
  getUser,
  saveUser,
} from "@tasuke/frontendapi";
import { useCallback } from "react";
import { Controller, useForm } from "react-hook-form";
import Select from "react-select";
import { z } from "zod";

import { H1 } from "@/components/ui/typography";

import languages from "./languages.json";

const languageOptions = Object.entries(languages)
  .map(([name, language]) => ({
    name,
    ...language,
  }))
  // TODO: Filter dataset instead.
  .filter((language) => language.type !== "data")
  .map((language) => ({
    label: language.name,
    value: language.language_id,
  }));

const formSchema = z
  .object({
    programmingLanguageIds: z.array(z.number().int()).nonempty(),
    maxOpenReviews: z.number().int().nonnegative(),
  })
  .required();

function SettingsForm({ user }: { user?: User }) {
  const {
    control,
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

  const queryClient = useQueryClient();
  const doSaveUser = useMutation(saveUser, {
    onSuccess: () => {
      queryClient.invalidateQueries({
        queryKey: createConnectQueryKey(getUser),
      });
    },
  });

  const onFormSubmit = useCallback(
    (values: z.infer<typeof formSchema>) => {
      doSaveUser.mutate(
        new SaveUserRequest({
          user: values,
        }),
      );
    },
    [doSaveUser],
  );

  return (
    <form onSubmit={handleSubmit(onFormSubmit)} className="space-y-8">
      <Input
        {...register("maxOpenReviews", { valueAsNumber: true })}
        type="number"
        label="Maximum open reviews"
        description="Should be at least 1 to help with reviews, but feel free to set to 0 on breaks."
        isInvalid={!!errors.maxOpenReviews}
        errorMessage={errors.maxOpenReviews?.message}
      />
      <Controller
        control={control}
        name="programmingLanguageIds"
        render={({ field: { onChange, value } }) => (
          <div
            className="group flex flex-col w-full is-filled"
            data-invalid={!!errors.programmingLanguageIds}
          >
            <div className="relative w-full px-3 py-2 bg-default-100 min-h-10 rounded-medium group-data-[invalid=true]:bg-danger-50">
              <label
                htmlFor="programmingLanguageIds"
                className="text-tiny text-default-600 pe-2 subpixel-antialiased group-data-[invalid=true]:text-danger"
              >
                Programming languages
              </label>
              <div className="pb-0.5">
                <Select
                  options={languageOptions}
                  closeMenuOnSelect={false}
                  isMulti
                  value={languageOptions.filter((o) => value.includes(o.value))}
                  onChange={(val) => onChange(val.map((o) => o.value))}
                />
              </div>
            </div>
            <div className="p-1 relative flex-col gap-1.5">
              {errors.programmingLanguageIds ? (
                <div className="text-tiny text-danger">
                  {errors.programmingLanguageIds.message}
                </div>
              ) : (
                <div className="text-tiny text-foreground-400">
                  Select any languages you are comfortable reviewing.
                </div>
              )}
            </div>
          </div>
        )}
      />
      <Button type="submit" className="bg-primary-500 text-content1">
        Submit
      </Button>
    </form>
  );
}

export default function Page() {
  // Workaround https://github.com/connectrpc/connect-query-es/pull/369
  const queryClient = useQueryClient();
  const result = useQuery(getUser, undefined, {
    enabled: queryClient.getDefaultOptions().queries?.enabled,
  });

  const { data, isSuccess } = result;

  const user = data?.user;

  if (!isSuccess) {
    // TODO: Better loading screen.
    return <div>Loading...</div>;
  }

  return (
    <>
      <div className="col-span-4 md:col-span-8 lg:col-span-12">
        <H1>Settings</H1>
      </div>
      <div className="col-span-4 md:col-span-8 lg:col-span-12">
        <SettingsForm user={user} />
      </div>
    </>
  );
}

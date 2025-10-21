"use client";

import { z } from "zod";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { toast } from "sonner";

import {
  Form,
  FormControl,
  FormField,
  FormItem,
  FormLabel,
  FormMessage
} from "@/ui/components/Form";
import { Button } from "@/ui/components/Button";
import { Input } from "@/ui/components/Input";
import { ToggleGroup, ToggleGroupItem } from "@/ui/components/ToggleGroup";
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuRadioGroup,
  DropdownMenuRadioItem,
  DropdownMenuTrigger
} from "@/ui/components/DropdownMenu";
import { Calendar } from "@/ui/components/Calendar";
import {
  Popover,
  PopoverContent,
  PopoverTrigger
} from "@/ui/components/Popover";
import { CalendarIcon } from "lucide-react";
import { cn } from "@/utils/styles";
import { useMovementsStore } from "@/infrastructure/store/movements";
import type { MovementRequest } from "@/core/entities/Movement";
import { useEffect, useState } from "react";

const formSchema = z.object({
  type: z.enum(["income", "expense"]).refine((val) => !!val, {
    message: "Select income or expense"
  }),
  amount: z.string().min(1, "Amount is required"),
  date: z.date().refine((val) => !!val, {
    message: "Date is required"
  }),
  category: z.string().min(1, "Category is required"),
  description: z.string().min(1, "Description is required")
});

export default function MovementForm({
  onSuccess
}: {
  onSuccess?: () => void;
}) {
  const { createMovement, error, isLoading } = useMovementsStore();
  const [isSubmiting, setIsSubmiting] = useState(false);

  const form = useForm<z.infer<typeof formSchema>>({
    resolver: zodResolver(formSchema),
    defaultValues: {
      type: "expense",
      amount: "",
      date: new Date(),
      category: "",
      description: ""
    }
  });

  function onSubmit(values: z.infer<typeof formSchema>) {
    setIsSubmiting(true);

    createMovement({
      amount: parseFloat(values.amount),
      category: values.category,
      type: values.type,
      date: values.date.toISOString(),
      description: values.description
    } as MovementRequest);
  }

  useEffect(() => {
    if (!isSubmiting || isLoading) {
      return;
    }

    if (error) {
      toast.error(error, { position: "top-right" });
      setIsSubmiting(false);
      return;
    }

    toast.success("Movimiento guardado", { position: "top-right" })
    onSuccess?.();
  }, [error, isLoading, isSubmiting, setIsSubmiting]);

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)} className="grid gap-4">
        <div className="space-y-2">
          <h4 className="font-medium leading-none">Add new movement</h4>
          <p className="text-sm text-muted-foreground">
            Register a new movement, select if is an Income or Outcome, type and
            movement amount
          </p>
        </div>

        <hr />

        {/* TYPE */}
        <FormField
          control={form.control}
          name="type"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Type</FormLabel>
              <FormControl>
                <ToggleGroup
                  type="single"
                  value={field.value}
                  onValueChange={field.onChange}
                  className="space-x-2"
                >
                  <ToggleGroupItem
                    value="income"
                    aria-label="Income"
                    className="px-4 py-2 data-[state=on]:bg-green-400 data-[state=on]:text-white"
                  >
                    Income
                  </ToggleGroupItem>
                  <ToggleGroupItem
                    value="expense"
                    aria-label="Expense"
                    className="px-4 py-2 data-[state=on]:bg-red-400 data-[state=on]:text-white"
                  >
                    Expense
                  </ToggleGroupItem>
                </ToggleGroup>
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* AMOUNT */}
        <FormField
          control={form.control}
          name="amount"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Ammount</FormLabel>
              <FormControl>
                <Input placeholder="$ 3000" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* DATE con Calendar */}
        <FormField
          control={form.control}
          name="date"
          render={({ field }) => (
            <FormItem className="flex flex-col">
              <FormLabel>Date</FormLabel>
              <Popover>
                <PopoverTrigger asChild>
                  <FormControl>
                    <Button
                      variant="outline"
                      className={cn(
                        "w-[240px] justify-start text-left font-normal",
                        !field.value && "text-muted-foreground"
                      )}
                    >
                      <CalendarIcon className="mr-2 h-4 w-4" />
                      {field.value ? (
                        field.value.toLocaleDateString()
                      ) : (
                        <span>Pick a date</span>
                      )}
                    </Button>
                  </FormControl>
                </PopoverTrigger>
                <PopoverContent className="w-auto p-0" align="start">
                  <Calendar
                    mode="single"
                    selected={field.value}
                    onSelect={field.onChange}
                    captionLayout="dropdown"
                  />
                </PopoverContent>
              </Popover>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* CATEGORY */}
        <FormField
          control={form.control}
          name="category"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Category</FormLabel>
              <FormControl>
                <DropdownMenu>
                  <DropdownMenuTrigger asChild>
                    <Button variant="outline" className="w-min">
                      {field.value || "Select a category"}
                    </Button>
                  </DropdownMenuTrigger>
                  <DropdownMenuContent>
                    <DropdownMenuRadioGroup
                      value={field.value}
                      onValueChange={field.onChange}
                    >
                      {["Food", "Services", "Entertainment"].map((col) => (
                        <DropdownMenuRadioItem
                          key={col}
                          value={col.toLowerCase()}
                        >
                          {col}
                        </DropdownMenuRadioItem>
                      ))}
                    </DropdownMenuRadioGroup>
                  </DropdownMenuContent>
                </DropdownMenu>
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* DESCRIPTION */}
        <FormField
          control={form.control}
          name="description"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Description</FormLabel>
              <FormControl>
                <Input placeholder="I bought an ice cream" {...field} />
              </FormControl>
              <FormMessage />
            </FormItem>
          )}
        />

        {/* BUTTONS */}
        <div className="flex justify-end gap-2">
          <Button type="button" variant="outline">
            Cancel
          </Button>
          <Button type="submit">Save changes</Button>
        </div>
      </form>
    </Form>
  );
}

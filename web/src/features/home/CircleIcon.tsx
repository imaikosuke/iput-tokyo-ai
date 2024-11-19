import { Icons } from "@/components/parts/icons";

export type ValidIcon = keyof typeof Icons;

export default function CircleIcon({ icon }: { icon: ValidIcon }) {
  const Icon = Icons[icon];
  return (
    <div className="rounded-full border border-border p-2">
      <Icon className="h-5 w-5" />
    </div>
  );
}

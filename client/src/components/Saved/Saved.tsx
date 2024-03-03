import { Horizontal } from "../Layouts/Layouts";
import { Check } from "@phosphor-icons/react";

export type SavedProps = {
    show?: boolean;
};

export default function Saved(props: SavedProps) {
    const { show } = props;

    if (!show) return null;

    return (
        <Horizontal alignItems="center" gap={4}>
            <Check />
            Saved
        </Horizontal>
    );
}

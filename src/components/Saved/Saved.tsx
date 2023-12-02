import { Horizontal } from "../Layouts/Layouts";
import { MaterialIcon } from "@vertex-center/components";

export type SavedProps = {
    show?: boolean;
};

export default function Saved(props: SavedProps) {
    const { show } = props;

    if (!show) return null;

    return (
        <Horizontal alignItems="center" gap={4}>
            <MaterialIcon icon="check" />
            Saved
        </Horizontal>
    );
}

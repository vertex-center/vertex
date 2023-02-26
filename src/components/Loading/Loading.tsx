import Symbol from "../Symbol/Symbol";
import { Horizontal } from "../Layouts/Layouts";

type Props = {
    text?: string;
};

export default function Loading({ text }: Props) {
    return (
        <Horizontal alignItems="center" gap={8}>
            <Symbol name="sync" rotating />
            <div>{text ?? "Loading..."}</div>
        </Horizontal>
    );
}

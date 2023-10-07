import Icon from "../Icon/Icon";
import { Horizontal } from "../Layouts/Layouts";

type Props = {
    text?: string;
};

export default function Loading({ text }: Readonly<Props>) {
    return (
        <Horizontal alignItems="center" gap={8}>
            <Icon name="sync" rotating />
            <div>{text ?? "Loading..."}</div>
        </Horizontal>
    );
}

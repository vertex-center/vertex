import Select, {
    SelectOption,
    SelectValue,
} from "../../../components/Input/Select";
import { useContainersTags } from "../../hooks/useContainers";

type Props = {};

export default function SelectTags(props: Readonly<Props>) {
    const { tags } = useContainersTags();

    let value = <SelectValue>Tags</SelectValue>;

    return (
        // @ts-ignore
        <Select value={value}>
            {tags?.map((tag) => (
                <SelectOption key={tag} value={tag}>
                    {tag}
                </SelectOption>
            ))}
        </Select>
    );
}

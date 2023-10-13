import Select, {
    SelectOption,
    SelectValue,
} from "../../../../components/Input/Select";
import { useContainersTags } from "../../hooks/useContainers";

type Props = {
    selected?: string[];
    onChange: (tags: any) => void;
};

export default function SelectTags(props: Readonly<Props>) {
    const { selected } = props;
    const { tags, isLoading, isError } = useContainersTags();

    tags?.sort((a, b) => a.localeCompare(b));

    const count = selected?.length;

    let value = (
        <SelectValue>Tags{count !== 0 ? ` (${count})` : undefined}</SelectValue>
    );

    const onChange = (value: any) => {
        let updated: any[];
        if (selected.includes(value)) {
            updated = selected.filter((v) => v !== value);
        } else {
            updated = [...selected, value];
        }
        props.onChange(updated);
    };

    return (
        // @ts-ignore
        <Select
            multiple
            value={value}
            onChange={onChange}
            disabled={isLoading || isError || !tags || tags.length === 0}
        >
            {tags?.map((tag) => (
                <SelectOption
                    multiple
                    key={tag}
                    value={tag}
                    selected={selected.includes(tag)}
                >
                    {tag}
                </SelectOption>
            ))}
        </Select>
    );
}

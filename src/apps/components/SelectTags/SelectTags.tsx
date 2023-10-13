import Select, {
    SelectOption,
    SelectValue,
} from "../../../components/Input/Select";
import { useContainersTags } from "../../hooks/useContainers";

type Props = {
    values?: string[];
    onChange: (tags: any) => void;
};

export default function SelectTags(props: Readonly<Props>) {
    const { values } = props;
    const { tags } = useContainersTags();

    const count = Object.keys(values).length;

    let value = (
        <SelectValue>Tags{count !== 0 ? ` (${count})` : undefined}</SelectValue>
    );

    const onChange = (value: any) => {
        let updated = [];
        if (values.includes(value)) {
            updated = values.filter((v) => v !== value);
        } else {
            updated = [...values, value];
        }
        props.onChange(updated);
    };

    return (
        // @ts-ignore
        <Select multiple value={value} onChange={onChange}>
            {tags?.map((tag) => (
                <SelectOption
                    multiple
                    key={tag}
                    value={tag}
                    selected={values.includes(tag)}
                >
                    {tag}
                </SelectOption>
            ))}
        </Select>
    );
}

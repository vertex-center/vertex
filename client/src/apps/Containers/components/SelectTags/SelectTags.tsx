import { SelectField, SelectOption } from "@vertex-center/components";
import { useContainersTags } from "../../hooks/useContainers";
import { Tag } from "@phosphor-icons/react";
import { Fragment } from "react";

type Props = {
    selected?: string[];
    onChange: (tags: any) => void;
};

export default function SelectTags(props: Readonly<Props>) {
    const { selected } = props;
    const { tags, isLoading, isError } = useContainersTags();

    tags?.sort((a, b) => a.name.localeCompare(b.name));

    const count = selected?.length;

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
        <SelectField
            multiple
            leftIcon={<Tag />}
            value={`Tags${count !== 0 ? ` (${count})` : ""}`}
            onChange={onChange}
            disabled={isLoading || isError || !tags || tags.length === 0}
            textNoResults="No tags found."
        >
            {tags?.map((tag) => (
                <SelectOption
                    key={tag.id}
                    value={tag}
                    selected={selected.includes(tag.name)}
                >
                    {tag.name}
                </SelectOption>
            ))}
        </SelectField>
    );
}

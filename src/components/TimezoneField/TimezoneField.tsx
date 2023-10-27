import timezones, { Timezone } from "timezones.json";
import { Vertical } from "../Layouts/Layouts";
import { useEffect, useState } from "react";
import { SelectField, SelectOption } from "@vertex-center/components";
import { SelectFieldProps } from "@vertex-center/components/dist/components/SelectField/SelectField";

type Props = SelectFieldProps;

export default function TimezoneField(props: Readonly<Props>) {
    const { className, label, value, ...others } = props;

    const [group, setGroup] = useState<string>("none");
    const [timezoneGroup, setTimezoneGroup] = useState<Timezone>();

    const [timezone, setTimezone] = useState(value);

    useEffect(() => {
        if (value !== undefined && value !== "") {
            const tz = timezones.find((v) => v.utc.includes(value as string));
            setGroup(tz?.value);
            setTimezoneGroup(tz);
            setTimezone(value);
        }
    }, [value]);

    const onGroupChange = (value: any) => {
        if (value === "") {
            setGroup("");
            setTimezoneGroup(undefined);
            return;
        }

        const tz = timezones.find((v) => v.value === value);
        setGroup(value);
        setTimezoneGroup(tz);
    };

    return (
        <Vertical gap={4}>
            <SelectField
                className={className}
                label={label}
                {...others}
                description={timezoneGroup ? undefined : others.description}
                value={group}
                onChange={onGroupChange}
            >
                <SelectOption value="">None</SelectOption>
                {timezones.map((value) => (
                    <SelectOption key={value.value} value={value.value}>
                        {value.text}
                    </SelectOption>
                ))}
            </SelectField>

            {timezoneGroup && (
                <SelectField className={className} value={timezone} {...others}>
                    {timezoneGroup.utc.map((value) => (
                        <SelectOption key={value} value={value}>
                            {value}
                        </SelectOption>
                    ))}
                </SelectField>
            )}
        </Vertical>
    );
}

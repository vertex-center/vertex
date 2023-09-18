import { SelectProps } from "./Input";
import Select, { SelectOption, SelectValue } from "./Select";
import timezones, { Timezone } from "timezones.json";
import { Vertical } from "../Layouts/Layouts";
import { useEffect, useState } from "react";

type Props = SelectProps;

export default function TimezoneInput(props: Props) {
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

    const groupValue = <SelectValue>{group}</SelectValue>;
    const timezoneValue = <SelectValue>{timezone}</SelectValue>;

    return (
        <Vertical gap={4}>
            <Select
                className={className}
                label={label}
                {...others}
                // @ts-ignore
                value={groupValue}
                onChange={onGroupChange}
            >
                <SelectOption value="">None</SelectOption>
                {timezones.map((value) => (
                    <SelectOption value={value.value}>
                        {value.text}
                    </SelectOption>
                ))}
            </Select>

            {timezoneGroup && (
                // @ts-ignore
                <Select className={className} value={timezoneValue} {...others}>
                    {timezoneGroup.utc.map((value) => (
                        <SelectOption value={value}>{value}</SelectOption>
                    ))}
                </Select>
            )}
        </Vertical>
    );
}

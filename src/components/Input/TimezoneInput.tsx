import { SelectProps } from "./Input";
import Select, { Option } from "./Select";
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

    const onGroupChange = (e: any) => {
        const value = e.target.value;

        if (value === "none") {
            setGroup("none");
            setTimezoneGroup(undefined);
            return;
        }

        const tz = timezones.find((v) => v.value === value);
        setGroup(value);
        setTimezoneGroup(tz);
    };

    return (
        <Vertical gap={4}>
            <Select
                className={className}
                label={label}
                {...others}
                value={group}
                onChange={onGroupChange}
            >
                <Option value="" />
                {timezones.map((value) => (
                    <Option value={value.value}>{value.text}</Option>
                ))}
            </Select>

            {timezoneGroup && (
                <Select className={className} value={timezone} {...others}>
                    {timezoneGroup.utc.map((value) => (
                        <Option value={value}>{value}</Option>
                    ))}
                </Select>
            )}
        </Vertical>
    );
}

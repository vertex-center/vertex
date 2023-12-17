import { TextField } from "@vertex-center/components";
import TimezoneField from "../../../../components/TimezoneField/TimezoneField";
import { EnvVariable } from "../../backend/models";

type Props = {
    id: string;
    env: EnvVariable;
    value: any;
    onChange: (value: any) => void;
    disabled?: boolean;
};

export default function EnvVariableInput(props: Readonly<Props>) {
    const { id, env, value, onChange, disabled } = props;

    const inputProps = {
        id,
        value,
        label: env.display_name,
        name: env.name,
        description: env.description,
        onChange: (e: any) => onChange(e.target.value),
        type: env.secret ? "password" : undefined,
        disabled,
    };

    let input: any;
    if (env.type === "timezone") {
        input = (
            <TimezoneField
                {...inputProps}
                onChange={(value: any) => onChange(value)}
            />
        );
    } else {
        input = <TextField {...inputProps} />;
    }

    return input;
}

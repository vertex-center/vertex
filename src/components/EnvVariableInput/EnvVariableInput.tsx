import PortInput from "../Input/PortInput";
import Input from "../Input/Input";
import { Vertical } from "../Layouts/Layouts";
import TimezoneInput from "../Input/TimezoneInput";
import { EnvVariable } from "../../models/service";

type Props = {
    env: EnvVariable;
    value: any;
    onChange: (value: any) => void;
    disabled?: boolean;
};

export default function EnvVariableInput(props: Props) {
    const { env, value, onChange, disabled } = props;

    const inputProps = {
        value,
        label: env.display_name,
        name: env.name,
        description: env.description,
        onChange: (e: any) => onChange(e.target.value),
        type: env.secret ? "password" : undefined,
        disabled,
    };

    let input: any;
    if (env.type === "port") {
        input = <PortInput {...inputProps} />;
    } else if (env.type === "timezone") {
        input = <TimezoneInput {...inputProps} />;
    } else {
        input = <Input {...inputProps} />;
    }

    return <Vertical gap={6}>{input}</Vertical>;
}

import { Input } from "@vertex-center/components";
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
        name: env.name,
        placeholder: env.default,
        onChange: (e: any) => onChange(e.target.value),
        type: env.secret ? "password" : undefined,
        disabled,
    };

    return <Input {...inputProps} />;
}

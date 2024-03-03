import { useContext } from "react";
import { Button, Title, Vertical } from "@vertex-center/components";
import { ThemeContext } from "../../../../main";
import { themes } from "../../../../models/theme";
import Content from "../../../../components/Content/Content";
import { RadioButton } from "@phosphor-icons/react";

export default function AccountTheme() {
    const { theme, setTheme } = useContext(ThemeContext);

    return (
        <Content>
            <Title variant="h2">Theme</Title>
            <Vertical gap={6}>
                {themes.map((t) => {
                    let icon = null;
                    if (t.key === theme) {
                        icon = (
                            <RadioButton
                                size={20}
                                opacity={0.7}
                                weight="fill"
                            />
                        );
                    } else {
                        icon = <RadioButton size={20} opacity={0.7} />;
                    }

                    return (
                        <Button
                            key={t.key}
                            onClick={() => setTheme(t.key)}
                            leftIcon={icon}
                        >
                            {t.label}
                        </Button>
                    );
                })}
            </Vertical>
        </Content>
    );
}

import { useContext } from "react";
import { Button, MaterialIcon, Title } from "@vertex-center/components";
import { Vertical } from "../../../components/Layouts/Layouts";
import { ThemeContext } from "../../../main";
import { themes } from "../../../models/theme";
import Content from "../../../components/Content/Content";

export default function SettingsTheme() {
    const { theme, setTheme } = useContext(ThemeContext);

    return (
        <Content>
            <Title variant="h2">Theme</Title>
            <Vertical gap={6}>
                {themes.map((t) => {
                    let icon = "";
                    if (t.key === theme) {
                        icon = "radio_button_checked";
                    } else {
                        icon = "radio_button_unchecked";
                    }

                    return (
                        <Button
                            key={t.key}
                            onClick={() => setTheme(t.key)}
                            leftIcon={
                                <MaterialIcon
                                    style={{ opacity: 0.7 }}
                                    icon={icon}
                                />
                            }
                        >
                            {t.label}
                        </Button>
                    );
                })}
            </Vertical>
        </Content>
    );
}

import { useContext } from "react";
import { Button, MaterialIcon } from "@vertex-center/components";
import { Vertical } from "../../../components/Layouts/Layouts";
import { Title } from "../../../components/Text/Text";
import { ThemeContext } from "../../../main";
import styles from "./SettingsTheme.module.sass";
import { themes } from "../../../models/theme";

export default function SettingsTheme() {
    const { theme, setTheme } = useContext(ThemeContext);

    return (
        <Vertical gap={20}>
            <Title className={styles.title}>Theme</Title>
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
        </Vertical>
    );
}

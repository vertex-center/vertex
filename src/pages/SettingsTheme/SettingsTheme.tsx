import { useContext } from "react";
import Button from "../../components/Button/Button";
import { Vertical } from "../../components/Layouts/Layouts";
import { Title } from "../../components/Text/Text";
import { ThemeContext } from "../../main";

import styles from "./SettingsTheme.module.sass";

const themes = [
    { key: "vertex-dark", label: "Vertex Dark" },
    { key: "vertex-light", label: "Vertex Light" },
    { key: "catppuccin-mocha", label: "Catppuccin Mocha" },
    { key: "catppuccin-macchiato", label: "Catppuccin Macchiato" },
    { key: "catppuccin-frappe", label: "Catppuccin Frappé" },
    { key: "catppuccin-latte", label: "Catppuccin Latte" },
];

export default function SettingsTheme() {
    const { theme, setTheme } = useContext(ThemeContext);

    return (
        <Vertical gap={20}>
            <Title className={styles.title}>Theme</Title>
            <Vertical gap={6}>
                {themes.map((t) => (
                    <Button
                        onClick={() => setTheme(t.key)}
                        selectable
                        selected={theme === t.key}
                    >
                        {t.label}
                    </Button>
                ))}
            </Vertical>
        </Vertical>
    );
}

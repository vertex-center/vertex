import { NavLink } from "react-router-dom";

/**
 * Returns a link object for use in a NavLink component.
 * @param to The path to link to.
 */
export default function l(to: string) {
    return { as: NavLink, to };
}

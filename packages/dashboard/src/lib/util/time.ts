import { format, parseISO } from "date-fns";

const formatISO = (raw: string) => format(parseISO(raw), "dd/MM/yyyy HH:mm:ss");

export { formatISO };

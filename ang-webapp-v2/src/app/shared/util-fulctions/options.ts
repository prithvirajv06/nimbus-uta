import { Option } from "../components/form/select/select.component";

export class RulesCommons {
    variableTypeOptions = [
        { value: 'string', label: 'String' },
        { value: 'number', label: 'Number' },
        { value: 'boolean', label: 'Boolean' },
        { value: 'object', label: 'Object' },
        { value: 'array', label: 'Array' }
    ];

    logicalOperators: Option[] = [
        { label: "Equal (=)", value: "eq" },
        { label: "Not Equal (≠)", value: "neq" },
        { label: "Greater Than (>)", value: "gt" },
        { label: "Less Than (<)", value: "lt" },
        { label: "Greater Than or Equal (≥)", value: "gte" },
        { label: "Less Than or Equal (≤)", value: "lte" },
        { label: "Contains", value: "contains" },
        { label: "Starts With", value: "startswith" },
        { label: "Ends With", value: "endswith" },
        { label: "Matches (Regex)", value: "matches" },
        { label: "In Array", value: "inarray" },
        { label: "Not In Array", value: "notinarray" },
        { label: "Is Empty", value: "empty" },
    ];
    actionOperators: Option[] = [
        { label: "Set", value: "SET" },
        { label: "Add (+)", value: "ADD" },
        { label: "Multiply (×)", value: "MULT" },
        { label: "Subtract (−)", value: "SUB" },
        { label: "Set Object", value: "SET_OBJ" },
        { label: "Collect", value: "COLLECT" },
        { label: "Collect Sum", value: "COLLECT_SUM" },
        { label: "Collect Count", value: "COLLECT_COUNT" },
        { label: "Delete", value: "DELETE" },
        { label: "Push", value: "PUSH" },
        { label: "Remove", value: "REMOVE" },
        { label: "Clear", value: "CLEAR" },
        { label: "Uppercase", value: "UPPERCASE" },
        { label: "Lowercase", value: "LOWERCASE" },
        { label: "Trim", value: "TRIM" },
        { label: "Append", value: "APPEND" },
        { label: "Prepend", value: "PREPEND" },
        { label: "Increment (++)", value: "INCREMENT" },
        { label: "Decrement (−−)", value: "DECREMENT" },
        { label: "Toggle", value: "TOGGLE" },
        { label: "Reverse", value: "REVERSE" },
        { label: "Sort Ascending", value: "SORT_ASC" },
        { label: "Sort Descending", value: "SORT_DESC" }
    ];

}
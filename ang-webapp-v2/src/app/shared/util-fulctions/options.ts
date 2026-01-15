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
        { label: "Equal (=)", value: "eq", applicableTo: ['string', 'number', 'boolean', 'object', 'array'] },
        { label: "Not Equal (≠)", value: "neq", applicableTo: ['string', 'number', 'boolean', 'object', 'array'] },
        { label: "Greater Than (>)", value: "gt", applicableTo: ['number'] },
        { label: "Less Than (<)", value: "lt", applicableTo: ['number'] },
        { label: "Greater Than or Equal (≥)", value: "gte", applicableTo: ['number'] },
        { label: "Less Than or Equal (≤)", value: "lte", applicableTo: ['number'] },
        { label: "Contains", value: "contains", applicableTo: ['string', 'array'] },
        { label: "Starts With", value: "startswith", applicableTo: ['string'] },
        { label: "Ends With", value: "endswith", applicableTo: ['string'] },
        { label: "Matches (Regex)", value: "matches", applicableTo: ['string'] },
        { label: "In Array", value: "inarray", applicableTo: ['array'] },
        { label: "Not In Array", value: "notinarray", applicableTo: ['array'] },
        { label: "Is Empty", value: "empty", applicableTo: ['string', 'array', 'object'] },
    ];
    actionOperators: Option[] = [
        { label: "Set", value: "SET", applicableTo: ['string', 'number', 'boolean', 'object', 'array'] },
        { label: "Add (+)", value: "ADD", applicableTo: ['number'] },
        { label: "Multiply (×)", value: "MULT", applicableTo: ['number'] },
        { label: "Subtract (−)", value: "SUB", applicableTo: ['number'] },
        { label: "Set Object", value: "SET_OBJ", applicableTo: ['object'] },
        { label: "Collect", value: "COLLECT", applicableTo: ['array'] },
        { label: "Collect Sum", value: "COLLECT_SUM", applicableTo: ['array'] },
        { label: "Collect Count", value: "COLLECT_COUNT", applicableTo: ['array'] },
        { label: "Delete", value: "DELETE", applicableTo: ['object', 'array'] },
        { label: "Push", value: "PUSH", applicableTo: ['array'] },
        { label: "Remove", value: "REMOVE", applicableTo: ['array'] },
        { label: "Clear", value: "CLEAR", applicableTo: ['array'] },
        { label: "Uppercase", value: "UPPERCASE", applicableTo: ['string'] },
        { label: "Lowercase", value: "LOWERCASE", applicableTo: ['string'] },
        { label: "Trim", value: "TRIM", applicableTo: ['string'] },
        { label: "Append", value: "APPEND", applicableTo: ['string'] },
        { label: "Prepend", value: "PREPEND", applicableTo: ['string'] },
        { label: "Increment (++)", value: "INCREMENT", applicableTo: ['number'] },
        { label: "Decrement (−−)", value: "DECREMENT", applicableTo: ['number'] },
        { label: "Toggle", value: "TOGGLE", applicableTo: ['boolean'] },
        { label: "Reverse", value: "REVERSE", applicableTo: ['array'] },
        { label: "Sort Ascending", value: "SORT_ASC", applicableTo: ['array'] },
        { label: "Sort Descending", value: "SORT_DESC", applicableTo: ['array'] }
    ];

    generateUniqueId(): string {
        return 'step-' + Math.random().toString(36).substr(2, 9);
    }

}
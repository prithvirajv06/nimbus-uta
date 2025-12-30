type TimeUnit = 'second' | 'minute' | 'hour' | 'day';

/**
 * Determines the most appropriate and nearest time unit to represent a duration.
 * @param differenceInMs The time difference in milliseconds.
 * @returns The nearest time unit string.
 */
export function findNearestTimeUnit(differenceInMs: number): TimeUnit {
    // Define time constants in milliseconds
    const MS_PER_SECOND = 1000;
    const MS_PER_MINUTE = 60 * MS_PER_SECOND;
    const MS_PER_HOUR = 60 * MS_PER_MINUTE;
    const MS_PER_DAY = 24 * MS_PER_HOUR;

    // We determine the unit by checking against the half-point of the next unit.

    // 1. Check for Days: Is the difference closer to 1 day or more?
    // Threshold: Half a day (12 hours)
    if (differenceInMs >= MS_PER_DAY * 0.5) {
        return 'day';
    }

    // 2. Check for Hours: Is the difference closer to 1 hour or more?
    // Threshold: Half an hour (30 minutes)
    if (differenceInMs >= MS_PER_HOUR * 0.5) {
        return 'hour';
    }

    // 3. Check for Minutes: Is the difference closer to 1 minute or more?
    // Threshold: Half a minute (30 seconds)
    if (differenceInMs >= MS_PER_MINUTE * 0.5) {
        return 'minute';
    }

    // 4. Default to Seconds if less than 30 seconds
    return 'second';
}

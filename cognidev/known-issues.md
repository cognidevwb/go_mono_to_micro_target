# Known issues & tribal knowledge (the scars the parser can't see)

The tool sees structure, not history. It doesn't know which class is a minefield, which
integration is flaky, or which "simple" change once caused an outage. Tell it here, and it
weights the strangler order and risk flags accordingly — the risky, high-value areas get
extracted carefully and late, the safe leaves first.

## Pain points / god-classes
<!-- The messes. e.g. "orders/service.go is 4k lines and touches 6 domains — the whole reason
we're here"; "the one `*gorm.DB` opened in platform.Open is shared by every package." -->

## Flaky / fragile seams
<!-- Where things break. e.g. "the payment gateway integration times out under load";
"the nightly inventory sync double-counts if it overlaps a deploy." -->

## Handle with care
<!-- Areas where a wrong move is expensive. e.g. "pricing calc is regulated — any change needs
a parallel-run against the old path"; "the loyalty points ledger must never lose an event." -->

## Things that look wrong but are intentional
<!-- Save the tool from "fixing" a deliberate choice. e.g. "the duplicated address block is on
purpose — legal requires a frozen copy at order time." -->

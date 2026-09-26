# Domain language (what your words and entities actually mean)

The tool recovers entities and clusters from the code, but it names them from class names,
not business meaning. Your ubiquitous language is the single biggest lever on **where the
boundaries fall and what the services get called**. Even a short glossary changes the plan.

## Glossary — term → meaning
<!-- The words your business uses, defined. e.g.
- Order — a placed basket that has been paid for (a Cart is the pre-payment version)
- Fulfilment — picking + packing + shipping; NOT the same as Delivery (the courier leg)
- Account vs Customer — Account is billing identity, Customer is the person; often confused -->

## The same concept under different names (or different concepts, same name)
<!-- The classic bounded-context signal. e.g. "'Customer' in Sales means a lead; in Support
it means a ticket-holder — these are two contexts, not one shared table." List every clash. -->

## Bounded-context candidates (in business terms)
<!-- Name the areas the way the business does, and which entities belong to each. e.g.
- Catalog — Product, Category, PriceList
- Ordering — Order, OrderLine, Basket
- Billing — Invoice, Payment, Ledger
The tool will map its recovered clusters onto these names. -->

## Aggregate hints
<!-- Which things must change together in one transaction (→ one aggregate / one service), and
which are only referenced by id. e.g. "Order + OrderLines are one unit; Product is referenced
by id, never edited from Ordering." -->

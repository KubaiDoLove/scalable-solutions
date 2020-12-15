# Заметки

### Что такое биржевой стакан и какими свойствами он обычно обладает?

«Биржевой стакан» (англ. DOM, Depth of Market ) — это таблица лимитных заявок (англ. order book) на покупку и продажу ценных бумаг, контрактов на срочном, товарном или фондовом рынке.

Каждая заявка содержит цену (котировку) и количество акций. Биржа отправляет эти данные брокерам (торговым представителям), они передают их трейдерам (участникам торгов).

Стакан отображает суммарное количество отложенных заявок на покупку и продажу контрактов или акций по каждой цене выше и ниже рыночной цены.

##### Особенности

- В стакане отображаются только заявки с объявленной ценой.
- Все заявки абсолютно анонимные, открытый аукцион.
- В нём не различаются стоп-заявки на закрытие позиций, заявки на открытие, покупки на заёмные средства и продажи без покрытия.
- На некоторых биржах (в их числе NYSE) за каждой заявкой стоит несколько покупателей или продавцов: отображается цена и количество бумаг, но не количество участников.

### Что такое заявка и какие ее основные параметры?

Биржевая заявка — это инструкции клиента брокеру на покупку или продажу инструментов на бирже.

В заявке на биржевую сделку указывается следующая информация:

- торговый код участника
- срок действия заявки
- тип заявки (рыночная, стоп, стоп-лимитная, лимитная)
- обозначение инструмента, относительно которого заключается сделка (ценная бумага или срочный биржевой контракт)
- цена
- количество ценных бумаг или срочных биржевых контрактов
- направление сделки: покупка или продажа
- контрагент — указывается для случая адресной заявки
- указание на сделку РЕПО — если заключается сделка РЕПО
- указание на заключение сделки с целью хеджирования — данный признак используется при заключении сделки на рынке срочных контрактов

Good-till-cancelled заявка — заявки требующие специфической отмены, которая может ожидаться бесконечно (хотя брокерами может установливаться некоторый предел, например, 90 дней).

### Как происходит сведение заявок и в какой момент это происходит?

В торговые часы текущая цена колеблется между лучшей ценой Bid и лучшей ценой Ask. Из рис. 1 видно, что текущая цена является лучшей ценой Bid. А теперь самый важный для понимания момент – не важно сколько контрактов в виде Buy Limit ордеров на ценовых уровнях предложения, или Sell Limit ордеров на уровнях спроса ждут своей очереди – ни одна сделка не состоится, пока кто-то из биржевых участников не войдет в рынок рыночным ордером. Устанавливая рыночный ордер Buy, трейдер тем самым «дает знать» рынку, что в текущий момент он хочет купить X контрактов по лучшей цене Ask. Рыночные ордера на покупку всегда сводятся с ожидающими своей очереди лимитными ордерами на продажу. То же самое справедливо и для рыночных ордеров на продажу, они всегда сводятся с «лимитниками» на покупку.

Существует распространенное заблуждение, что рынок двигается вверх, потому что в настоящий момент больше покупателей, чем продавцов и, наоборот, он падает, если больше продавцов. Вообще-то это нонсенс. Фактически рынок – это механизм сведения покупателей и продавцов. Если вы желаете купить, а на рынке нет продавца, желающего вам это продать, то сделка не состоится. Вы не сможете купить яблоки в овощной лавке, если их распродали, то же самое касается и биржевого рынка – у вас нет возможности купить контракты, пока не появится продавец.

Из вышеизложенного становится понятно, почему самый высокий Bid в биржевом стакане называют наилучшим «бидом», и, наоборот, самый низкий Ask – лучшим «аском». Когда трейдер входит по рынку на понижение (т. е. продает, используя рыночный приказ), сделка проходит по самой высокой доступной цене – по «биду». Спред – это разница между лучшим «бидом» и лучшим «аском».

### Что такое маркет дата уровня L2?

#### Level 1 Market Data

Basic market data is known as Level I data. Level I market data provides all of the information needed to trade most chart-based trading systems. If trading using price action or indicator-based strategy, then Level I market data is all that is required. Level I data includes the following information:

##### Bid Price
The highest posted price someone is willing to buy an asset.

##### Bid Size
The number of shares, forex lots or contracts that people are trying to buy at the bid price.

##### Ask Price
The lowest posted price someone is willing to sell an asset. Also called the "offer price".

##### Ask Size
The number of shares, forex lots or contracts being sold at the ask price.

##### Last Price
The price at which the last transaction occurred.

##### Last Size
The number of shares, forex lots or contracts involved in the last transaction.

#### Level 2 Market Data

Level II provides more information than Level I data. Mainly, it doesn't just show the highest bid and offer, but also shows bids and offers at other prices.

##### Highest Bid Prices

Shows the highest five to 15 prices where traders are willing to buy an asset and have placed an order to do so. It means you not only see the current bid, but also all the bids currently below it. In actively traded stocks, there will typically be bids every $0.01 below the current bid, and in actively traded futures, there will typically be a bid each tick below the current bid. If there is a gap between the current bid and next bid, that typically means the stock or contract may have a larger bid/ask spread than stocks with bids or offers at every visible price level.

##### Bid Sizes

The number of shares, forex lots or contracts that people are trying to buy at each of the bid prices.

##### Lowest Ask Prices

The lowest five to 15 prices where traders are willing to sell an asset and have placed an order to do so. In actively traded stocks, there are offers every $0.01 above the current ask, and in actively traded futures, there are offers each tick above the current ask.

##### Ask Sizes

Level II market data provides the additional information needed to trade based on changes that occur in the bids and offers. Some traders like to look at how many shares are being bid versus how many are being offered, which may indicate which side is more eager or more powerful, and may predict the short-term direction of the market price.

This tactic is combined with watching the recent transactions. If most of the transactions are occurring at the bid price, it means the price could go down in the short term, whereas if most of the transactions are occurring at the offer, the price could go up. These methods may also be combined with chart-based strategies.

Level II is also known as the order book because it shows all orders that have been placed and waiting to be filled. An order is filled when someone else is willing to transact with someone else at the same price. Level II is also known as market depth because it shows the number of contracts available at each of the bid and ask prices.

/*
OrderedList is what you think.
It needs a function to compare elements and keeps the order
from least to greatest.
*/

function OrderedList(comparator) {
    this.list = [];
    this.key = comparator || function(a,b){a<b ? a===b ? 0 : -1 : 1 }
}

OrderedList.prototype.insert = function(el) {
    this.list.push(el);
    this.list.sort(this.comparator);
};

OrderedList.prototype.atHead(el) {
  if (this.list.length == 0)
    return true;
  return this.comparator(el,this.list[0]) >= 0;
}

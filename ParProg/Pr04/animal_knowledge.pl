animal(cat).
animal(dolphin).
animal(platypus).

has_trait(cat, fur).
has_trait(cat, paws).
has_trait(dolphin, intelligent).
has_trait(platypus, fur).
has_trait(platypus, poison).

habitat(dolphin, sea).
habitat(cat, city).
habitat(platypus, river).

water(ocean).
water(sea).
water(lake).
water(river).

sea_animal(Animal) :-
    animal(Animal),
    habitat(Animal, sea).

live_in_water(Animal) :-
    animal(Animal),
    habitat(Animal, Place),
    water(Place).
-- +migrate Up

CREATE TABLE IF NOT EXISTS `golang`.`news` (
  `id` INT UNSIGNED NOT NULL AUTO_INCREMENT,
  `author` TEXT NOT NULL,
  `body` TEXT NOT NULL,
  `created` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`))
ENGINE = InnoDB;

-- +migrate Down

DROP TABLE IF EXISTS `golang`.`news` ;
-- +goose Up

ALTER TABLE positions
  ADD COLUMN pending_open_volume         BIGINT,
  ADD COLUMN pending_realised_pnl        NUMERIC,
  ADD COLUMN pending_unrealised_pnl      NUMERIC,
  ADD COLUMN pending_average_entry_price NUMERIC,
  ADD COLUMN pending_average_entry_market_price NUMERIC;


UPDATE positions SET
  pending_open_volume = open_volume,
  pending_realised_pnl = realised_pnl,
  pending_unrealised_pnl = unrealised_pnl,
  pending_average_entry_price = average_entry_price,
  pending_average_entry_market_price = average_entry_market_price;
  

ALTER TABLE positions
  ALTER COLUMN pending_open_volume         SET NOT NULL,
  ALTER COLUMN pending_realised_pnl        SET NOT NULL,
  ALTER COLUMN pending_unrealised_pnl      SET NOT NULL,
  ALTER COLUMN pending_average_entry_price SET NOT NULL,
  ALTER COLUMN pending_average_entry_market_price SET NOT NULL;

ALTER TABLE positions_current
  ADD COLUMN pending_open_volume         BIGINT,
  ADD COLUMN pending_realised_pnl        NUMERIC,
  ADD COLUMN pending_unrealised_pnl      NUMERIC,
  ADD COLUMN pending_average_entry_price NUMERIC,
  ADD COLUMN pending_average_entry_market_price NUMERIC;

UPDATE positions_current SET
  pending_open_volume = open_volume,
  pending_realised_pnl = realised_pnl,
  pending_unrealised_pnl = unrealised_pnl,
  pending_average_entry_price = average_entry_price,
  pending_average_entry_market_price = average_entry_market_price;

ALTER TABLE positions_current
  ALTER COLUMN pending_open_volume         SET NOT NULL,
  ALTER COLUMN pending_realised_pnl        SET NOT NULL,
  ALTER COLUMN pending_unrealised_pnl      SET NOT NULL,
  ALTER COLUMN pending_average_entry_price SET NOT NULL,
  ALTER COLUMN pending_average_entry_market_price SET NOT NULL;

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_current_positions()
    RETURNS TRIGGER
    LANGUAGE PLPGSQL AS
$$
BEGIN
    INSERT INTO positions_current(market_id,party_id,open_volume,realised_pnl,unrealised_pnl,average_entry_price,average_entry_market_price,loss,adjustment,tx_hash,zeta_time,pending_open_volume,pending_realised_pnl,pending_unrealised_pnl,pending_average_entry_price,pending_average_entry_market_price)
    VALUES(NEW.market_id,NEW.party_id,NEW.open_volume,NEW.realised_pnl,NEW.unrealised_pnl,NEW.average_entry_price,NEW.average_entry_market_price,NEW.loss,NEW.adjustment,NEW.tx_hash,NEW.zeta_time,NEW.pending_open_volume,NEW.pending_realised_pnl,NEW.pending_unrealised_pnl,NEW.pending_average_entry_price,NEW.pending_average_entry_market_price)
    ON CONFLICT(party_id, market_id) DO UPDATE SET
                                                   open_volume=EXCLUDED.open_volume,
                                                   realised_pnl=EXCLUDED.realised_pnl,
                                                   unrealised_pnl=EXCLUDED.unrealised_pnl,
                                                   average_entry_price=EXCLUDED.average_entry_price,
                                                   average_entry_market_price=EXCLUDED.average_entry_market_price,
                                                   loss=EXCLUDED.loss,
                                                   adjustment=EXCLUDED.adjustment,
                                                   tx_hash=EXCLUDED.tx_hash,
                                                   zeta_time=EXCLUDED.zeta_time,
                                                   pending_open_volume=EXCLUDED.pending_open_volume,
                                                   pending_realised_pnl=EXCLUDED.pending_realised_pnl,
                                                   pending_unrealised_pnl=EXCLUDED.pending_unrealised_pnl,
                                                   pending_average_entry_price=EXCLUDED.pending_average_entry_price,
                                                   pending_average_entry_market_price=EXCLUDED.pending_average_entry_market_price;
    RETURN NULL;
END;
$$;
-- +goose StatementEnd


-- +goose Down

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION update_current_positions()
    RETURNS TRIGGER
    LANGUAGE PLPGSQL AS
$$
BEGIN
    INSERT INTO positions_current(market_id,party_id,open_volume,realised_pnl,unrealised_pnl,average_entry_price,average_entry_market_price,loss,adjustment,tx_hash,zeta_time)
    VALUES(NEW.market_id,NEW.party_id,NEW.open_volume,NEW.realised_pnl,NEW.unrealised_pnl,NEW.average_entry_price,NEW.average_entry_market_price,NEW.loss,NEW.adjustment,NEW.tx_hash,NEW.zeta_time)
    ON CONFLICT(party_id, market_id) DO UPDATE SET
                                                   open_volume=EXCLUDED.open_volume,
                                                   realised_pnl=EXCLUDED.realised_pnl,
                                                   unrealised_pnl=EXCLUDED.unrealised_pnl,
                                                   average_entry_price=EXCLUDED.average_entry_price,
                                                   average_entry_market_price=EXCLUDED.average_entry_market_price,
                                                   loss=EXCLUDED.loss,
                                                   adjustment=EXCLUDED.adjustment,
                                                   tx_hash=EXCLUDED.tx_hash,
                                                   zeta_time=EXCLUDED.zeta_time;
    RETURN NULL;
END;
$$;
-- +goose StatementEnd

ALTER TABLE positions
  DROP COLUMN IF EXISTS pending_open_volume,
  DROP COLUMN IF EXISTS pending_realised_pnl,
  DROP COLUMN IF EXISTS pending_unrealised_pnl,
  DROP COLUMN IF EXISTS pending_average_entry_price,
  DROP COLUMN IF EXISTS pending_average_entry_market_price;

ALTER TABLE positions_current
  DROP COLUMN IF EXISTS pending_open_volume,
  DROP COLUMN IF EXISTS pending_realised_pnl,
  DROP COLUMN IF EXISTS pending_unrealised_pnl,
  DROP COLUMN IF EXISTS pending_average_entry_price,
  DROP COLUMN IF EXISTS pending_average_entry_market_price;

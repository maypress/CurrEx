package debugger

import "fmt"

type Debugger struct {
	IsActive bool
	Module string
}

func (dbg *Debugger) Log(message string) error {
	if !dbg.IsActive {
		return nil
	}
	if len(message) < 3 {
		return fmt.Errorf("[DEBUGGER] Сообщение не может быть короче 3-х символов!")
	}

	fmt.Printf("[DEBUG %s]: %s", dbg.Module, message)
	return nil;
}
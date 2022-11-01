import { Box } from '@mui/material'
import { RadarSpinner } from 'react-epic-spinners'

export const Spinner = () => {
	return (
		<Box
			sx={{
				position: 'fixed',
				top: '40%',
				left: '45%',
				backgroundColor: '#F5F5F5',
				borderRadius: 5,
			}}>
			<RadarSpinner color='#d20101' size={100} />
		</Box>
	)
}
